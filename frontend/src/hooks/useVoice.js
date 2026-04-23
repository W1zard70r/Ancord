import { useState, useRef, useCallback } from 'react';

export const useVoice = (userId, chatId) => {
    const [status, setStatus] = useState('disconnected');
    const [remoteStreams, setRemoteStreams] = useState([]);

    const pcRef = useRef(null);
    const wsRef = useRef(null);
    const localStreamRef = useRef(null);

    const joinVoice = useCallback(async () => {
        if (!userId || !chatId) return;
        setStatus('connecting...');

        try {
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
            localStreamRef.current = stream;

            const turnRes = await fetch('http://localhost:8081/turn-config');
            const turnConfig = await turnRes.json();

            const pc = new RTCPeerConnection(turnConfig);
            pcRef.current = pc;

            pc.onicecandidate = (event) => {
                if (event.candidate && wsRef.current?.readyState === WebSocket.OPEN) {
                    wsRef.current.send(JSON.stringify({ type: 'candidate', payload: event.candidate.candidate }));
                }
            };

            pc.ontrack = (event) => {
                const newStream = event.streams[0];
                setRemoteStreams(prev => {
                    if (prev.find(s => s.id === newStream.id)) return prev;
                    return [...prev, newStream];
                });
            };

            stream.getTracks().forEach(track => pc.addTrack(track, stream));

            const ws = new WebSocket(`ws://127.0.0.1:8081/ws?user_id=${userId}&chat_id=${chatId}`);
            wsRef.current = ws;

            ws.onopen = async () => {
                const offer = await pc.createOffer();
                await pc.setLocalDescription(offer);
                ws.send(JSON.stringify({ type: 'offer', payload: offer.sdp }));
            };

            ws.onmessage = async (event) => {
                const msg = JSON.parse(event.data);
                if (msg.type === 'answer') {
                    await pc.setRemoteDescription({ type: 'answer', sdp: msg.payload });
                    setStatus('connected');
                } else if (msg.type === 'renegotiate') {
                    const offer = await pc.createOffer();
                    await pc.setLocalDescription(offer);
                    ws.send(JSON.stringify({ type: 'offer', payload: offer.sdp }));
                }
            };

            ws.onclose = () => leaveVoice();

        } catch (err) {
            console.error("Voice Error:", err);
            setStatus('error');
        }
    }, [userId, chatId]);

    const leaveVoice = useCallback(() => {
        if (wsRef.current) wsRef.current.close();
        if (pcRef.current) pcRef.current.close();
        if (localStreamRef.current) {
            localStreamRef.current.getTracks().forEach(track => track.stop());
        }

        wsRef.current = null;
        pcRef.current = null;
        localStreamRef.current = null;

        setRemoteStreams([]);
        setStatus('disconnected');
    }, []);

    return { status, joinVoice, leaveVoice, remoteStreams };
};
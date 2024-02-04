import {useEffect, useState} from 'react';

export default function App() {
    const [data, setData] = useState("");
    const getData = async () => {
        const resp = await fetch('https://805d-2405-9800-bca0-7cce-6021-33b3-ca1e-e4a5.ngrok-free.app/course');
        const json = await resp.json();
        setData(json);
        console.log(json);
    }
    
    

    useEffect(() => {
        getData();
    }, []);

    return (
        <pre>
        {JSON.stringify(data, null, 2)}
        </pre>
    )
}
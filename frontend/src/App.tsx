import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, GetCurrentTime} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [timeText, setTimeText] = useState("Click to get current time");
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    function getCurrentTime() {
        GetCurrentTime().then(setTimeText);
    }

    return (
        <div id="App">
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div>
            <div id="time-section" className="input-box">
                <div id="time-result" className="result">{timeText}</div>
                <button className="btn" onClick={getCurrentTime}>GetCurrentTime</button>
            </div>
        </div>
    )
}

export default App

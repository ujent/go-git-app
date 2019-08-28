import React from 'react';
import '../../App.css';

const Log = props => {

    return (
        <li>
            <button type="button" className="button" disabled={!props.isAvailable} onClick={() => props.action()}>Log</button>
        </li>
    );
}

export default Log;
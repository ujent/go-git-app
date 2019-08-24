import React from 'react';
import '../../App.css';

const Pull = props => {
    return (
        <li>
            <div className="command-block command-block-content">
                <button type="button" className="button">Pull</button>
                <input type="text" placeholder="remote" className="command-block-input" />
            </div>
            <div className="credentials">
                <input type="text" placeholder="user" />
                <input type="password" placeholder="password" />
            </div>
            <hr></hr>
        </li>
    );
}

export default Pull;
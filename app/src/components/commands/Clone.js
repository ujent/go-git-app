import React from 'react';
import '../../App.css';

const Clone = props => {
    return (
        <li>
            <div className="command-block command-block-content">
                <button type="button" className="button">Clone</button>
                <input type="text" placeholder="URL" className="command-block-input" />
            </div>
            <div className="credentials">
                <input type="text" placeholder="user" />
                <input type="password" placeholder="password" />
            </div>
            <hr></hr>
        </li>
    );
};

export default Clone;
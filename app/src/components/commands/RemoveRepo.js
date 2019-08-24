import React from 'react';
import '../../App.css';

const RemoveRepo = props => {
    return (
        <li className="command-block">
            <div className="command-block-content">
                <button type="button" className="button medium-button">Remove repo</button>
                <select placeholder="select" className="command-block-input" defaultValue="">
                    <option value="" disabled hidden>select</option>
                    <option value="repo1">repo1</option>
                    <option value="repo2">repo2</option>
                </select>
            </div>
            <hr></hr>
        </li>
    );
}

export default RemoveRepo;
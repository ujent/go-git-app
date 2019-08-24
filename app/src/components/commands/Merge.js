import React from 'react';
import '../../App.css';

const Merge = props => {

    return (
        <li className="command-block">
            <div className="command-block-content">
                <button type="button" className="button">Merge</button>
                <select placeholder="branch" className="command-block-input" defaultValue="">
                    <option value="" disabled hidden>select</option>
                    <option value="branch1">branch1</option>
                    <option value="branch2">branch2</option>
                </select></div>

            <hr></hr>
        </li>
    );
};

export default Merge;
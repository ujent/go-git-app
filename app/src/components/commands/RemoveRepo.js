import React from 'react';
import '../../App.css';

const RemoveRepo = props => {

    const repoOptions = props.repositories.map(
        opt => {
            return (
                <option key={opt} value={opt}>{opt}</option>
            );
        }
    );

    return (
        <li className="command-block">
            <div className="command-block-content">
                <button type="button" className="button medium-button" disabled={!props.isAvailable}>Remove repo</button>
                <select placeholder="select" className="command-block-input" defaultValue="">
                    <option value="" disabled hidden>select</option>
                    {repoOptions}
                </select>
            </div>
            <hr></hr>
        </li>
    );
}

export default RemoveRepo;
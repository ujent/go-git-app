import React from 'react';
import '../../App.css';

const Merge = props => {

    const branchOptions = props.branches.map(
        opt => {
            return (
                <option key={opt} value={opt}>{opt}</option>
            );
        }
    );

    return (
        <li className="command-block">
            <div className="command-block-content">
                <button type="button" className="button" disabled={!props.isAvailable}>Merge</button>
                <select placeholder="branch" className="command-block-input" defaultValue="">
                    <option value="" disabled hidden>select</option>
                    {branchOptions}
                </select></div>

            <hr></hr>
        </li>
    );
};

export default Merge;
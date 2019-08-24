import React from 'react';
import '../../App.css';

const Commit = props => {

    return (
        <li className="commit-command">
            <button type="button" className="button" disabled={!props.isAvailable}>Commit</button>
            <textarea placeholder="commit message" rows="3"></textarea>
            <hr></hr>
        </li>
    );
}

export default Commit;
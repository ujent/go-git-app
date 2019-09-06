import React from 'react';
import '../App.css';
import spinnerIcon from '../icons/spinner.svg'

const Spinner = props => {
    if (!props.isVisible) {
        return null;
    }

    return (
        <div className="spinner-overlay">
            <div className="spinner-wrapper">
                <img src={spinnerIcon} alt="spinner" />
            </div>
        </div>
    );
}

export default Spinner;
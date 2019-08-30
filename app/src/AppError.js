import React from 'react';
import './App.css';

const AppError = props => {
    return (
        <div className="app-wrapper">

            <div className="app-main">
                <div>{props.error}</div>
            </div>
        </div>
    );
}

export default AppError;
import React from 'react';
import '../App.css';


const ConfirmPopup = props => {
    const { isOpen, message, onClose, onConfirm, closePopup } = props;
    if (!isOpen) {
        return null;
    }

    function closeWithCancel() {
        if (onClose && typeof (onClose) === 'function') {
            onClose();
        }

        closePopup();
    }

    function closeWithConfirm() {
        if (onConfirm && typeof (onConfirm) === 'function') {
            onConfirm();
        }

        closePopup();
    }

    return (
        <div className="popup-overlay" onClick={closeWithCancel}>
            <section className="popup-wrapper">
                <header className="popup-header">Confirmation</header>
                <div className="popup" onClick={e => e.stopPropagation()}>
                    <div className="popup-message">
                        <p>{message}</p>
                    </div>
                    <div className="popup-buttons-wrapper">
                        <button className="button" onClick={closeWithCancel}>Cancel</button>
                        <button className="button" onClick={closeWithConfirm}>OK</button>
                    </div>
                </div>
            </section>
        </div>
    );
};

export default ConfirmPopup;
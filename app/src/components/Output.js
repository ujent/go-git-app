import React from 'react';

const Output = props => {
    const { message } = props;

    return (
        <section className="output">
            <textarea disabled rows="6" placeholder="output" value={message}></textarea>
        </section>
    );
}

export default Output;
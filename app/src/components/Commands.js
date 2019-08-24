import React from 'react';
import '../App.css';
import Commit from "./commands/Commit"
import Merge from "./commands/Merge"
import Clone from "./commands/Clone"
import Pull from "./commands/Pull"
import Push from "./commands/Push"
import RemoveRepo from "./commands/RemoveRepo"
import CheckoutBranch from "./commands/CheckoutBranch"
import RemoveBranch from "./commands/RemoveBranch"
import Log from "./commands/Log"

const Commands = props => {

    return (
        <section className="commands">
            <h2>Commands</h2>
            <ul className="commands-list">
                <Commit />
                <Merge />
                <Clone />
                <Pull />
                <Push />
                <RemoveRepo />
                <CheckoutBranch />
                <RemoveBranch />
                <Log />

            </ul>
        </section>
    );
};

export default Commands;
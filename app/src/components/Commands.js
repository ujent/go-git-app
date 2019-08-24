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
                <Commit action={props.handleCommit} />
                <Merge action={props.handleMerge} branches={props.branches} />
                <Clone action={props.handleClone} />
                <Pull action={props.handlePull} />
                <Push action={props.handlePush} />
                <RemoveRepo action={props.handleRemoveRepo} repositories={props.repositories} />
                <CheckoutBranch action={props.handleCheckoutBranch} branches={props.branches} />
                <RemoveBranch action={props.handleRemoveBranch} branches={props.branches} />
                <Log action={props.handleLog} />
            </ul>
        </section>
    );
};

export default Commands;
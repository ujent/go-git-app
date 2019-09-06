import React from 'react';
import '../App.css';
import Commit from "./commands/Commit"
import Merge from "./commands/Merge"
import Clone from "./commands/Clone"
import Pull from "./commands/Pull"
import Push from "./commands/Push"
import RemoveRepo from "./commands/RemoveRepo"
import CreateBranch from "./commands/CreateBranch"
// import CheckoutBranch from "./commands/CheckoutBranch"
import RemoveBranch from "./commands/RemoveBranch"
import Log from "./commands/Log"

const Commands = props => {
    const { repositories, branches, currentUser, currentRepo, currentBranch } = props;
    const isBaseAvailable = currentUser !== "" && currentRepo !== "" && currentBranch !== "";
    const canCommit = currentUser !== "" && currentRepo !== "";
    const canCheckoutBranch = currentUser !== "" && currentRepo !== "";
    const canRemoveBranch = currentUser !== "" && currentRepo !== "";

    return (
        <section className="commands">
            <h2>Commands</h2>
            <ul className="commands-list">
                <Commit action={props.handleCommit} isAvailable={canCommit} />
                <Merge mergeAction={props.handleMerge} abortAction={props.handleAbortMerge} branches={branches} currentBranch={currentBranch} isAvailable={isBaseAvailable} />
                <Clone action={props.handleClone} isAvailable={currentUser !== ""} />
                <Pull action={props.handlePull} isAvailable={isBaseAvailable} />
                <Push action={props.handlePush} isAvailable={isBaseAvailable} />
                <RemoveRepo action={props.handleRemoveRepo} repositories={repositories} isAvailable={currentUser !== ""} />
                {/* <CheckoutBranch action={props.handleCheckoutBranch} branches={branches} currentBranch={currentBranch} isAvailable={canCheckoutBranch} /> */}
                <CreateBranch action={props.handleCheckoutBranch} branches={branches} isAvailable={canCheckoutBranch} />
                <RemoveBranch action={props.handleRemoveBranch} branches={branches} isAvailable={canRemoveBranch} />
                <Log action={props.handleLog} isAvailable={isBaseAvailable} />
            </ul>
        </section>
    );
};

export default Commands;
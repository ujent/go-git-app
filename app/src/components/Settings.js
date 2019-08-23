import React from 'react';
import '../App.css';

const Settings = props => {
    const { users, selectedUser, repos, selectedRepo, branches, selectedBranch } = props;
    const userOptions = users.map(
        opt => {
            return (
                <option key={opt} value={opt}>{opt}</option>
            );
        }
    );

    const repoOptions = repos.map(
        opt => {
            return (
                <option key={opt} value={opt}>{opt}</option>
            );
        }
    );

    const isReposDisabled = repos.length === 0;

    const branchOptions = branches.map(
        opt => {
            return (
                <option key={opt} value={opt}>{opt}</option>
            );
        }
    );

    const isBranchesDisabled = branches.length === 0;

    return (
        <section className="main-settings">
            <div>
                <label htmlFor="userSelectID">User</label>
                <select id="userSelectID" value={selectedUser} onChange={(e) => props.handleSelectUser(e.target.value)}>
                    <option value="" disabled hidden>select</option>
                    {userOptions}
                </select>
            </div>
            <div>
                <label htmlFor="repoSelectID">Repository</label>
                <select id="repoSelectID" disabled={isReposDisabled} value={selectedRepo} onChange={(e) => props.handleSelectRepo(e.target.value)}>
                    <option value="" disabled hidden>select</option>
                    {repoOptions}
                </select>
            </div>
            <div>
                <label htmlFor="branchSelectID">Branch</label>
                <select id="branchSelectID" disabled={isBranchesDisabled} value={selectedBranch} onChange={(e) => props.handleSelectBranch(e.target.value)}>
                    <option value="" disabled hidden>select</option>
                    {branchOptions}
                </select>
            </div>
        </section>
    );
};

export default Settings;

import React, { Component } from 'react';
import '../../App.css';

export default class Merge extends Component {

    constructor(props) {
        super(props);

        this.state = {
            branch: ''
        }
    }



    render() {
        const isAvailable = this.props.isAvailable && this.state.branch
        const branchOptions = this.props.branches.map(
            opt => {
                if (opt !== this.props.currentBranch) {
                    return (
                        <option key={opt} value={opt}>{opt}</option>
                    );
                } else {
                    return '';
                }
            }
        );

        return (
            <li className="command-block">
                <div className="command-block-content">
                    <button type="button" className="button" disabled={!isAvailable} onClick={this.onMergeClick} >Merge</button>
                    <select placeholder="branch" className="command-block-input" value={this.state.branch} onChange={this.onBranchChange}>
                        <option value="" disabled hidden>select</option>
                        {branchOptions}
                    </select></div>
                <hr></hr>
            </li>
        );
    };

    onBranchChange = (e) => {
        this.setState({
            branch: e.target.value
        })
    }

    onMergeClick = () => {
        this.props.action(this.state.branch);
        this.setState({
            branch: ''
        })
    }


};

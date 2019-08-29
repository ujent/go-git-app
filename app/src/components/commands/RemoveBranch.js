import React, { Component } from 'react';
import '../../App.css';

export default class RemoveBranch extends Component {
    constructor(props) {
        super(props);

        this.state = {
            branch: ''
        }
    }

    render() {

        const isAvailable = this.props.isAvailable && this.state.branch;

        const branchOptions = this.props.branches.map(
            opt => {
                return (
                    <option key={opt} value={opt}>{opt}</option>
                );
            }
        );

        const onBranchChange = (e) => {
            this.setState({
                branch: e.target.value
            });
        }

        const onRemoveClick = () => {
            this.props.action(this.state.branch)

            this.setState({
                branch: ''
            });
        }

        return (
            <li className="command-block">
                <div className="command-block-content">
                    <button type="button" className="button medium-button" disabled={!isAvailable} onClick={onRemoveClick} >Remove branch</button>
                    <select placeholder="select" className="command-block-input" value={this.state.branch} onChange={onBranchChange} >
                        <option value="" disabled hidden>select</option>
                        {branchOptions}
                    </select>
                </div>
                <hr></hr>
            </li>
        );
    }
}

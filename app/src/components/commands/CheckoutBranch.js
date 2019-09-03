import React, { Component } from 'react';
import '../../App.css';

export default class CheckoutBranch extends Component {
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
                if (opt === this.props.currentBranch) {
                    return '';
                }

                return (
                    <option key={opt} value={opt}>{opt}</option>
                );
            }
        );

        return (
            <li className="command-block">
                <div className="command-block-content">
                    <button type="button" className="button medium-button" disabled={!isAvailable} onClick={this.onCheckoutClick} >Checkout branch</button>
                    <select placeholder="select" className="command-block-input" value={this.state.branch} onChange={this.onBranchChange} >
                        <option value="" disabled hidden>select</option>
                        {branchOptions}
                    </select>
                </div>
                <hr></hr>
            </li>
        );
    }

    onCheckoutClick = () => {
        this.props.action(this.state.branch);

        this.setState({
            branch: ''
        });
    }

    onBranchChange = (e) => {
        this.setState({
            branch: e.target.value
        });
    }


}

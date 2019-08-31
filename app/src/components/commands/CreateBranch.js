import React, { Component } from 'react';
import '../../App.css';

export default class CreateBranch extends Component {
    constructor(props) {
        super(props);

        this.state = {
            branch: ''
        }
    }

    render() {
        const isAvailable = this.props.isAvailable && this.state.branch;

        const onCreateClick = () => {
            this.props.action(this.state.branch);

            this.setState({
                branch: ''
            });
        }

        const onBranchChange = (e) => {
            this.setState({
                branch: e.target.value
            });
        }

        return (
            <li className="command-block">
                <div className="command-block-content">
                    <button type="button" className="button medium-button" disabled={!isAvailable} onClick={onCreateClick} >Create branch</button>
                    <input type="text" placeholder="branch" className="command-block-input" value={this.state.branch} onChange={onBranchChange} />
                </div>
                <hr></hr>
            </li>
        );
    }

}

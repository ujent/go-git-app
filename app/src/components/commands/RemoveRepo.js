import React, { Component } from 'react';
import '../../App.css';

export default class RemoveRepo extends Component {
    constructor(props) {
        super(props);

        this.state = {
            repo: '',
        }
    }
    render() {
        const isAvailable = this.props.isAvailable && this.state.repo;

        const repoOptions = this.props.repositories.map(
            opt => {
                return (
                    <option key={opt} value={opt}>{opt}</option>
                );
            }
        );

        return (
            <li className="command-block">
                <div className="command-block-content">
                    <button type="button" className="button medium-button" disabled={!isAvailable} onClick={this.onRemoveClick} >Remove repo</button>
                    <select placeholder="select" className="command-block-input" value={this.state.repo} onChange={this.onRepoChange} >
                        <option value="" disabled hidden>select</option>
                        {repoOptions}
                    </select>
                </div>
                <hr></hr>
            </li>
        );
    };

    onRepoChange = (e) => {
        this.setState({
            repo: e.target.value
        });
    }

    onRemoveClick = () => {
        this.props.action(this.state.repo)

        this.setState({
            repo: ''
        });
    }


}

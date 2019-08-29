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

        const onRepoChange = (e) => {
            this.setState({
                repo: e.target.value
            });
        }

        const onRemoveClick = () => {
            this.props.action(this.state.repo)

            this.setState({
                repo: ''
            });
        }

        return (
            <li className="command-block">
                <div className="command-block-content">
                    <button type="button" className="button medium-button" disabled={!isAvailable} onClick={onRemoveClick} >Remove repo</button>
                    <select placeholder="select" className="command-block-input" value={this.state.repo} onChange={onRepoChange} >
                        <option value="" disabled hidden>select</option>
                        {repoOptions}
                    </select>
                </div>
                <hr></hr>
            </li>
        );
    };

}

import React, { Component } from 'react';
import '../../App.css';

export default class Push extends Component {
    constructor(props) {
        super(props);

        this.state = {
            url: '',
            name: '',
            psw: ''
        }
    }

    render() {
        const onURLChange = (e) => {
            this.setState({
                url: e.target.value
            });
        }

        const onAuthNameChange = (e) => {
            this.setState({
                name: e.target.value
            });
        }

        const onAuthPswChange = (e) => {
            this.setState({
                psw: e.target.value
            });
        }

        const onPushClick = () => {
            this.props.action(this.state.url, this.state.name, this.state.psw);
            this.setState({
                url: '',
                name: '',
                psw: ''
            });
        }

        return (
            <li>
                <div className="command-block command-block-content">
                    <button type="button" className="button" disabled={!this.props.isAvailable} onClick={onPushClick} >Push</button>
                    <input type="text" placeholder="remote" className="command-block-input" value={this.state.url} onChange={onURLChange} />
                </div>
                <div className="credentials">
                    <input type="text" placeholder="user" value={this.state.name} onChange={onAuthNameChange} />
                    <input type="password" placeholder="password" value={this.state.psw} onChange={onAuthPswChange} />
                </div>
                <hr></hr>
            </li>
        );
    };

}

import React, { Component } from 'react';
import '../../App.css';

export default class Commit extends Component {
    constructor(props) {
        super(props);

        this.state = {
            comment: ''
        }
    }

    render() {

        return (
            <li className="commit-command">
                <button type="button" className="button" disabled={!this.props.isAvailable} onClick={() => this.onCommitClick()}>Commit</button>
                <textarea placeholder="commit message" rows="3" value={this.state.comment} onChange={(e) => this.onCommentChange(e)}></textarea>
                <hr></hr>
            </li>
        );
    }

    onCommentChange = (e) => {
        this.setState({
            comment: e.target.value
        });
    }

    onCommitClick = () => {
        this.props.action(this.state.comment);

        this.setState({
            comment: ''
        })
    }

}

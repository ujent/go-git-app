import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import './index.css';
import App from './containers/AppContainer';
import * as serviceWorker from './serviceWorker';
import createAppStore from './createAppStore';
import { switchUser, getRepositories, getBranches, getFiles } from './api';
import { StorageItem } from './constants';


// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();


const currentUser = window.localStorage.getItem(StorageItem.User);
const currentRepo = window.localStorage.getItem(StorageItem.Repo);
const currentBranch = window.localStorage.getItem(StorageItem.Branch);

if (currentUser !== null) {
    switchUser(currentUser).then(
        () => {
            getRepositories(currentUser).then(
                rs => {
                    if (rs.repos.indexOf(currentRepo) !== -1) {
                        getBranches(currentUser, currentRepo).then(
                            brRS => {
                                if (brRS.branches.indexOf(currentBranch) !== -1) {

                                    getFiles(currentUser, currentRepo, currentBranch).then(
                                        filesRS => {
                                            const store = createAppStore(null, currentUser, currentRepo, currentBranch, rs.repos, brRS.branches, filesRS.files);

                                            ReactDOM.render(
                                                <Provider store={store}><App /></Provider>,
                                                document.getElementById('root')
                                            );
                                        },
                                        err => {
                                            console.log(err)
                                            const store = createAppStore(err, currentUser, currentRepo, currentBranch, rs.repos, brRS.branches, []);

                                            ReactDOM.render(
                                                <Provider store={store}><App /></Provider>,
                                                document.getElementById('root')
                                            );
                                        }
                                    );

                                } else {
                                    const store = createAppStore(null, currentUser, currentRepo, '', rs.repos, brRS.branches);
                                    window.localStorage.removeItem(StorageItem.Branch)

                                    ReactDOM.render(
                                        <Provider store={store}><App /></Provider>,
                                        document.getElementById('root')
                                    );
                                }


                            },
                            err => {
                                console.log(err)
                                const store = createAppStore(err, currentUser, currentRepo, '', rs.repos);

                                ReactDOM.render(
                                    <Provider store={store}><App /></Provider>,
                                    document.getElementById('root')
                                );
                            }
                        )
                    } else {
                        const store = createAppStore(null, currentUser, '', '', rs.repos);
                        window.localStorage.removeItem(StorageItem.Repo)

                        ReactDOM.render(
                            <Provider store={store}><App /></Provider>,
                            document.getElementById('root')
                        );
                    }
                },
                err => {
                    console.log(err)
                    const store = createAppStore(err, currentUser);

                    ReactDOM.render(
                        <Provider store={store}><App /></Provider>,
                        document.getElementById('root')
                    );
                }
            )
        },
        err => {
            console.log(err)
            const store = createAppStore(err);

            ReactDOM.render(
                <Provider store={store}><App /></Provider>,
                document.getElementById('root')
            );
        }
    )
} else {
    const store = createAppStore();

    ReactDOM.render(
        <Provider store={store}><App /></Provider>,
        document.getElementById('root')
    );
}
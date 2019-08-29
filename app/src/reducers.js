import { ActionType } from './constants';

export const rootReducer = (state = {}, action) => {
    switch (action.type) {
        case ActionType.SET_CURRENT_USER: {
            const prev = state.settings.currentUser;
            if (prev === action.user) {
                return state;
            }

            return Object.assign({}, state, {
                repositories: [],
                branches: [],
                settings: Object.assign({}, state.settings, {
                    currentUser: action.user,
                    currentRepo: '',
                    currentBranch: ''
                }),
                files: [],
                fileContent: {
                    isVisible: false,
                    content: ''
                }
            });
        }
        case ActionType.SET_BRANCHES: {
            return Object.assign({}, state, {
                branches: action.branches
            });
        }
        case ActionType.SET_CURRENT_BRANCH: {
            return Object.assign({}, state, {
                settings: Object.assign({}, state.settings, {
                    currentBranch: action.current,
                })
            });
        }
        case ActionType.REMOVE_BRANCH: {
            const currentBranch = state.settings.currentBranch;
            let branches = state.branches;

            for (let i = 0; i < branches.length; i++) {
                if (branches[i] === action.branch) {
                    branches.splice(i, 1);
                    break
                }
            }

            if (currentBranch === action.branch) {
                return Object.assign({}, state, {
                    branches: branches,
                    settings: Object.assign({}, state.settings, {
                        currentBranch: ''
                    }),
                    files: [],
                    fileContent: {
                        isVisible: false,
                        content: ''
                    }
                });
            }

            return Object.assign({}, state, {
                branches: branches
            });
        }
        case ActionType.REMOVE_REPO: {
            const currentRepo = state.settings.currentRepo;
            let repos = state.repositories;

            for (let i = 0; i < repos.length; i++) {
                if (repos[i] === action.repo) {
                    repos.splice(i, 1);
                }

            }

            if (currentRepo === action.repo) {
                return Object.assign({}, state, {
                    repositories: repos,
                    settings: Object.assign({}, state.settings, {
                        currentRepo: ''
                    }),
                    fileContent: {
                        isVisible: false,
                        content: ''
                    },
                    files: []
                })
            }

            return Object.assign({}, state, {
                repositories: repos,
            })
        }
        case ActionType.SET_REPOSITORIES: {
            return Object.assign({}, state, {
                repositories: action.repos,
            });
        }
        case ActionType.SET_CURRENT_REPOSITORY: {
            return Object.assign({}, state, {
                settings: Object.assign({}, state.settings, {
                    currentRepo: action.current,
                })
            });
        }
        case ActionType.SHOW_MSG: {
            return Object.assign({}, state, {
                output: action.msg
            })
        }
        case ActionType.RESET_MSG: {
            return Object.assign({}, state, {
                output: ''
            })
        }
        default:
            return state;
    }
};
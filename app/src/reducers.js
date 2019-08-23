import { ActionType } from './constants';

export const rootReducer = (state = {}, action) => {
    switch (action.type) {
        case ActionType.SET_USERS: {
            const prev = state.settings.currentRepo
            return Object.assign({}, state, {
                users: action.users,
                settings: Object.assign({}, state.settings, {
                    currentUser: action.current,
                    prevUser: prev
                })
            });
        }
        case ActionType.SET_BRANCHES: {
            const prev = state.settings.currentRepo
            return Object.assign({}, state, {
                branches: action.branches,
                settings: Object.assign({}, state.settings, {
                    currentBranch: action.current,
                    prevBranch: prev
                })
            });
        }
        case ActionType.RESET_BRANCH: {
            const prev = state.settings.currentBranch
            return Object.assign({}, state, {
                settings: Object.assign({}, state.settings, {
                    currentBranch: "",
                    prevBranch: prev
                })
            });
        }
        case ActionType.SET_REPOSITORIES: {
            const prev = state.settings.currentRepo
            return Object.assign({}, state, {
                repositories: action.repos,
                settings: Object.assign({}, state.settings, {
                    currentRepo: action.current,
                    prevRepo: prev
                })
            });
        }
        case ActionType.RESET_REPO: {
            const prev = state.settings.currentRepo
            return Object.assign({}, state, {
                settings: Object.assign({}, state.settings, {
                    currentRepo: "",
                    prevRepo: prev
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
                output: ""
            })
        }
        default:
            return state;
    }
};
import { ActionType } from './constants';

export const rootReducer = (state = {}, action) => {
    switch (action.type) {
        // case 'ActionType.CHANGE_PROJECTS_LIST_SEARCH_FOCUS':
        //   return Object.assign({}, state, {
        //     projectsListPage: Object.assign(state.projectsListPage, {
        //       filters: Object.assign(state.projectsListPage.filters, {
        //         isSearchFocused: !state.projectsListPage.filters.isSearchFocused
        //       })
        //     })
        //   });

        default:
            return state;
    }
};
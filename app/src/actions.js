import * as api from './api';
import { ActionType } from './constants';

export function changeProjectsListSearchFocus() {
    return {
        //type: ActionType.CHANGE_PROJECTS_LIST_SEARCH_FOCUS
    };
}
/* export function changeBenefitsFilter(newFilter) {
  return {
    type: ActionType.CHANGE_BENEFITS_FILTER,
    newFilter
  };
}

export function startFreeco() {
  return (dispatch, getState) => {
    const user = getUser(getState);
    api.startFreeco(user).then(
      () => {
        dispatch(getProcessInfo());
      },
      err => {
        dispatch(showError(err));
      }
    );
  };
}*/
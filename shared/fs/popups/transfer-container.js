// @flow
import * as Types from '../../constants/types/fs'
import * as Constants from '../../constants/fs'
import * as FsGen from '../../actions/fs-gen'
import {compose, connect, setDisplayName, type TypedState, type Dispatch} from '../../util/container'
import {navigateUp} from '../../actions/route-tree'
import TransferPopup from './transfer'
import {formatDurationFromNowTo} from '../../util/timestamp'

const mapStateToProps = (state: TypedState, {routeProps}) => {
  const path = routeProps.get('path')
  const _pathItem = state.fs.pathItems.get(path) || Constants.makeUnknownPathItem()
  const _username = state.config.username || undefined
  const _transfers = state.fs.transfers
  return {
    path,
    _pathItem,
    _username,
    _transfers,
  }
}

const mapDispatchToProps = (dispatch: Dispatch) => ({
  _onHidden: () => dispatch(navigateUp()),
  _dismissTransfer: (key: string) => dispatch(FsGen.createDismissTransfer({key})),
})

const mergeProps = (stateProps, {_onHidden, _dismissTransfer}) => {
  const itemStyles = Constants.getItemStyles(
    Types.getPathElements(stateProps.path),
    stateProps._pathItem.type,
    stateProps._username
  )
  const [key, {completePortion, endEstimate, isDone, intent}] =
    stateProps._transfers
      .filter(ts => ts.type === 'download' && ts.intent !== 'none')
      .entries()
      .next().value || Constants.makeTransferState()
  const onHidden = () => {
    _onHidden()
    _dismissTransfer(key)
  }
  isDone && onHidden()
  return {
    name: stateProps._pathItem.name,
    intent,
    itemStyles,
    completePortion,
    progressText: formatDurationFromNowTo(endEstimate),
    onHidden,
  }
}

export default compose(
  connect(mapStateToProps, mapDispatchToProps, mergeProps),
  setDisplayName('TransferPopup')
)(TransferPopup)
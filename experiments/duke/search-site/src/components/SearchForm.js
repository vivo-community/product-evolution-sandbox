import React, { Component } from 'react'
import PropTypes from 'prop-types'
import SearchField from './SearchField'
import { requestSearch } from '../actions/search'

export class SearchForm extends Component {

  constructor(props, context) {
    super(props, context)
    this.handleSubmitSearch = this.handleSubmitSearch.bind(this)
  }

  handleSubmitSearch(e) {
    e.preventDefault()
    
    const { search : { searchFields }, dispatch } = this.props
    const term = this.term.value
    let start = 0

    const constructedSearch = { start: start, term: term }
    console.debug("trying to search")
    console.debug(constructedSearch)
    
    dispatch(requestSearch(constructedSearch))

  }

  render() {
    const { search : { isFetching, searchFields } } = this.props
    let query = { ...searchFields }
    
    let button = <button type="submit" className="btn btn-primary btn-sm">Search</button>

    return (
        <form onSubmit={this.handleSubmitSearch} className="form-horizontal">
            <SearchField ref={(ref) => this.term = ref} defaultValue={query.term} autofocus={true}/>
            {button}
        </form>
    )
  }

}

import { connect } from 'react-redux'
//import { bindActionCreators } from 'redux'

const mapStateToProps = (search, ownProps) => {
  return search
}

//function mapDispatchToProps(dispatch) {
//  return bindActionCreators({search}, dispatch)
//}

export default connect(mapStateToProps)(SearchForm)

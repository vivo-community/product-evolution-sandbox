import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Page from '../layouts/page'
import SearchForm from '../components/SearchForm'
import SearchResults from '../components/SearchResults'

export class VivoSearchApp extends Component {

  constructor(props, context) {
    super(props, context)
  }

  componentDidMount() {
    const { dispatch } = this.props
  }

  render() {
    return (
      <Page>
        <SearchForm />
        <hr />
        <SearchResults /> 
      </Page>
    )
  }

}


import { connect } from 'react-redux'

const mapStateToProps = (search, ownProps) => {
  return { ...search }
}

export default connect(mapStateToProps)(VivoSearchApp)

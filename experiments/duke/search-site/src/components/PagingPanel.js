import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import ReactDOM from 'react-dom'
import classNames from 'classnames'
import { requestSearch } from '../actions/search'

import * as constants from '../actions/types'

import helper from '../pager'

import PagingLink from './PagingLink'
import PageFlipper from './PageFlipper'
import _ from 'lodash'

export class PagingPanel extends Component {

  constructor(props, context) {
    super(props, context)
    this.handlePage = this.handlePage.bind(this) 
  }

  handlePage(e, pageNumber) {
    e.preventDefault()
    // given page - what is start ??
    const { search : { searchFields }, dispatch } = this.props

    // given page (in parameter) calculate start
    let start = searchFields ? searchFields['start'] : 0
    let newStart = (pageNumber - 1) * constants.PAGE_ROWS

    // NOTE: if not a new 'query' obj (like below) - this error happens:
    // useQueries.js:35 Uncaught TypeError: object.hasOwnProperty is not a function
    const query = { ...searchFields, start: newStart }

    let full_query = { ...query }

    dispatch(requestSearch(full_query))
    
    /*
    this.context.router.push({
      pathname: '/',
      query: full_query
    })
    */
  
  }
  
  render() {
    // so start should be coming from search object (state)
    const { search : { results, searchFields } } = this.props

    let { hits } = results
    let { total = 0 } = hits

    let start = searchFields['start'] || 0

    if (!hits) {
      return ( <div></div> )
    }

    console.debug("paging")
    console.debug(`start=${start}`)
    console.debug(`total=${total}`)

    // 105 results
    // start at 50
    // would be page 2 of 3
    // NOTE: all these Math.floor(s) are annoying
    console.debug(`${total}/${constants.PAGE_ROWS}`)

    let totalPages = Math.floor(total/constants.PAGE_ROWS)
    const remainder = total % constants.PAGE_ROWS
    
    if (remainder) { totalPages +=1 }

    if (totalPages == 0) {
      return ( <div></div> )
    }

    const currentPage = Math.floor(start/constants.PAGE_ROWS) + 1
 
    console.debug(`totalPages=${totalPages} currentPage=${currentPage}`)

    let pageMap = helper.pageArrays(totalPages, currentPage)
    // pageMap is an array set of arrays
    // more/less links are returned as ['+', 16] or ['-'] (means no number)
    //
    // so example might be [['+', 1][16...30]['+', 31]]
 
    let [previous, current, next] = pageMap
         
    let pages = _.map(current, (x) => {
      let active = (x == currentPage) ? true : false
      let key = `pageLinkTo_${x}`

      console.debug("trying to make page link")
      return (
          <PagingLink key={key} pageNumber={x} active={active} onClick={(e) => this.handlePage(e, x)}/>
      )
    })

    console.debug(`previous=${previous} from ${pageMap}`)

    //let backward = flip(previous, 'backward') 
    //let forward = flip(next, 'forward') 
          /*
          <ul className="pagination">
            <PageFlipper pages={previous} direction="backward" onClick={(e) => this.handlePage(e, pageNumber)} />
            {pages}
            <PageFlipper pages={next} direction="forward" onClick={(e) => this.handlePage(e, pageNumber)} />
          </ul>
          */
    return (
        <nav>
          <ul className="pagination">
            <PageFlipper pages={previous} direction="backward" onClick={(e) => this.handlePage(e, pageNumber)} />
            {pages}
            <PageFlipper pages={next} direction="forward" onClick={(e) => this.handlePage(e, pageNumber)} />
          </ul>
        </nav>
      )
  }
}

const mapStateToProps = (search, ownProps) => {
  return  search;
}

export default connect(mapStateToProps)(PagingPanel);

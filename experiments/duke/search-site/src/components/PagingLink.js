import React, { Component } from 'react'
import PropTypes from 'prop-types'
import ReactDOM from 'react-dom'

export default class PagingLink extends Component {

  constructor(props, context) {
    super(props, context)
    this.onClick = this.props.onClick
    //this.onClick = this.onClick.bind(this)
  }
  
  /*
  onClick(e) {
    e.preventDefault() 
  }
  */

  renderPage(pageNumber, active) {
      let key = `pageLinkTo_${pageNumber}`

      if(active) {
        return (
         <li className="page-item active">
           <span className="page-link">{pageNumber}</span>
         </li>
        )
      } else {
         return (
          <li className="page-item">
            <span className="page-link">
              <a href="#" onClick={(e) => this.onClick(e, pageNumber)}>{pageNumber}</a>
            </span>
          </li>
        )
      }
  }
  
  render() {
    return this.renderPage(this.props.pageNumber, this.props.active)
  }
}




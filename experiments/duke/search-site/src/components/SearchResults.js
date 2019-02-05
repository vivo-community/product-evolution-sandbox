import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import ReactDOM from 'react-dom'

import _ from 'lodash'
import PagingPanel from './PagingPanel'

//require('../styles/vivo_search.less');

export class GenericResult extends Component {
  
  render() { 
    let { id, uri, label } = this.props.thing
    return (
       <div key={id} className={`panel`}>
         <a href={uri}>{ label }</a>
         <div><span className={`badge badge-secondary`}>{this.props.type}</span></div>
       </div>
      )
    } 
}

export class GrantResult extends Component {
  
  render() { 
    let { id, uri, label } = this.props.thing
    return (
       <div key={id} className={`panel`}>
         <a href={uri}>{ label }</a>
         <div><span className={`badge badge-info`}>{this.props.type}</span></div>
       </div>
      )
    } 
}


export class PersonResult extends Component {
  
  render() { 
    let { id, uri, label } = this.props.thing
    let { name } = this.props.thing
    return (
       <div key={id} className={`panel`}>
         <a href={uri}>{ name.firstName} { name.lastName }</a>
         <div><span className={`badge badge-success`}>{this.props.type}</span></div>
       </div>
      )
    } 
}

export class PublicationResult extends Component {
  
  render() { 
    let { id, uri, label } = this.props.thing
    return (
       <div key={id} className={`panel`}>
         <a href={uri}>{ label }</a>
         <div><span className={`badge badge-primary`}>{this.props.type}</span></div>
       </div>
      )
    } 
}


export class SearchResult extends Component {
  
  constructor(props, context) {
    super(props, context)
  }

  render() { 
    let { id, uri, label } = this.props.thing

    switch(this.props.type) {
      case 'person':
        return (<PersonResult key={id} thing={this.props.thing} type={this.props.type} />)
      case 'publication':
        return (<PublicationResult key={id} thing={this.props.thing } type={this.props.type} />)
      case 'grant':
        return (<GrantResult key={id} thing={this.props.thing } type={this.props.type} />)
      default:
        return (<GenericResult key={id} thing={this.props.thing } type={this.props.type} />)
    }
    /*
    if (this.props.type == 'person') { 
      let { name } = this.props.thing
      return (
       <div key={id}>
         <a href={uri}>{ id }:{ name.lastName }</a>
         <span>{this.props.type}</span>
       </div>
      )
    } else {
     return (
       <div key={id}>
         <a href={uri}>{ id }:{ label }</a>
         <span>{this.props.type}</span>
       </div>
     )
    }
    */
  }

}

export class SearchResults extends Component {

  constructor(props, context) {
    super(props, context)
  }

  shouldComponentUpdate(nextProps, nextState) {
    console.debug("should update search results?")
    console.debug(nextProps)
    return true
  }

  render() {
    const { search : { results } } = this.props
 
    console.log("trying to render results:")
    console.debug(results)
    let { hits = 0 } = results

    console.debug(hits)
    const items = []
    let spanClass = "badge-success"
    if (hits) {
        let x = hits.hits
        _.forEach(x, function(value, key) {
          let source = value._source
          items.push(<SearchResult key={source.id} type={value._type} thing={ source } />)
        })
    }
    if (hits) {
       return (
        <section className="search-results">
          <h2>Search Results</h2>
          <div>
            { items }
          </div>
          <hr />
          <PagingPanel />
        </section>
       )
    } else {
      return (
        <section className="search-results">
          <h2>Search Results</h2>
          <p>No results to show</p>
        </section>
      )
    }
    
  }

}


const mapStateToProps = (search, ownProps) => {
  return { ...search }
}

export default connect(mapStateToProps)(SearchResults)

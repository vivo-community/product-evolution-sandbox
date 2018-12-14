import React from 'react'
import { withRouteData } from 'react-static'
import { Link } from '@reach/router'
import PersonTile from './person_tile'
import './people.css'

const People = ({people, currentPage, totalPages}) => {
  return (
    <div className="department-content">
      <h1>People</h1>
      { totalPages ?
        <div className="pager">
          { [...Array(totalPages).keys()].map((pg) => {
            let pageNumber = pg + 1
            let isCurrentPage = pageNumber === currentPage
            return <div key={`people_page_${pageNumber}`} className="page">
              { isCurrentPage ?
                <b>{pageNumber}</b>
              : <Link to={`/people/page/${pageNumber}`}>{pageNumber}</Link>
              }
            </div>
          })}
        </div>
      : null }
      <div className="people">
        { people.map((p) => <PersonTile key={p.id} person={p}/>)}
      </div>
    </div>
  )
}

export default withRouteData(People)


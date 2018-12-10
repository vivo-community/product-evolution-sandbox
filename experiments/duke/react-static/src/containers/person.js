import React from 'react'
import { withRouteData, Link } from 'react-static'
import DisplayDate from './display_date.js'
import DateRange from './date_range.js'
import './person.css'

const Person = ({person}) => {
  let {
    id,
    name: {
      firstName,
      middleName,
      lastName
    },
    image: {
      main
    },
    affiliationList,
    overviewList,
    educationList,
    publicationList,
    grantList
  } = person

  let displayName = [firstName,middleName,lastName].filter((n) => n).join(" ")
  let displayTitle = affiliationList[0].label
  return (
    <div className="person">
      <div className="person-header">
        <img className="person-image" src={main} alt={`${displayName} Profile Image`}/>
        <h1 className="person-name">{displayName}</h1>
        <div className="person-title">{displayTitle}</div>
        { overviewList.map(({overview, type: { code }}) => {
          return (
            <div
              key={code}
              className="person-overview"
              dangerouslySetInnerHTML={{__html: overview}}
            />
          )
        }) }
        { affiliationList.length > 0 ?
        <div className="person-collection">
          <h3>Current Appointments and Affiliations</h3>
          { affiliationList.map((affiliation) => {
            let {dateTime, resolution} = affiliation.startDate
            return (
              <div key={affiliation.id} className="person-affiliation person-collection-item">
                <span className="affiliation-label">{affiliation.label}</span>
                <span className="affiliation-date">
                  <DisplayDate dateTime={dateTime} resolution={resolution}/>
                </span>
              </div>
            )
          })}
        </div>
        : null }
        { educationList.length > 0 ?
        <div className="person-collection">
          <h3>Education</h3>
          { educationList.map((education) => {
            return (
              <div key={education.id} className="person-education person-collection-item">
                <span className="education-label">{education.label}</span>
                <span className="education-org">, {education.org.label}</span>
              </div>
            )
          })}
        </div>
        : null }
      </div>
      { publicationList.length > 0 ?
      <div className="person-collection">
        <h3>Publications</h3>
        { publicationList.map((publication) => {
          return (
            <div key={publication.id} className="person-publication person-collection-item">
              { publication.authorList ?
              <span className="publication-authors">{publication.authorList}. </span>
              : null }
              <span className="publication-label">{publication.label}</span>
              { publication.venue.label ?
              <span className="publication-venue">. {publication.venue.label}</span>
              : null }
            </div>
          )
        })}
      </div>
      : null }
      { grantList.length > 0 ?
      <div className="person-collection">
        <h3>Grants</h3>
        { grantList.map((grant) => {
          return (
            <div key={grant.id} className="person-grant person-collection-item">
              <span className="grant-label">{grant.label}</span>
              <span className="grant-daterange">
                <DateRange startDate={grant.startDate} endDate={grant.endDate}/>
              </span>
            </div>
          )
        })}
      </div>
      : null }
    </div>
  )
  return <p>person</p>

}

export default withRouteData(Person)

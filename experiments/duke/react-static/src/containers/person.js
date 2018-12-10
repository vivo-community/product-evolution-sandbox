import React from 'react'
import { withRouteData, Link } from 'react-static'
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
    publicationList
  } = person

  let displayName = [firstName,middleName,lastName].filter((n) => n).join(" ")
  let displayTitle = affiliationList[0].label
  return (
    <div className="person">
      <div className="person-header">
        <img className="person-image" src={main} alt={`${displayName} Profile Image`}/>
        <div className="person-name">{displayName}</div>
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
      </div>
      { affiliationList.length > 0 ?
      <div className="person-collection">
        <h3>Current Appointments and Affiliations</h3>
        { affiliationList.map((affiliation) => {
          let {dateTime, resolution} = affiliation.startDate
          let displayDate
          if (resolution === 'year') {
            displayDate = new Date(dateTime).getFullYear();
          }
          return (
            <div key={affiliation.id} className="person-affiliation person-collection-item">
              <span className="affiliation-label">{affiliation.label}</span>
              <span className="affiliation-date">{displayDate}</span>
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
    </div>
  )
  return <p>person</p>

}

export default withRouteData(Person)

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
    overviewList
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
    </div>
  )
  return <p>person</p>

}

export default withRouteData(Person)

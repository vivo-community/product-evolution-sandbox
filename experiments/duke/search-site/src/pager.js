import _ from 'lodash'

/*
 * NOTE: this returns an array of 3 arrays given a total number of pages
 * and the current page.  The first array is what to do with *before*,
 * the last array is what to do with *after*
 *
 * so, as an example (assuming PAGE_BY = 15):
 *
 * if we have 95 pages, and we 
 * are on page 1:
 *
 [ [ '-' ],
 [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15 ],
 [ '+', 16 ] ]
 - means no page to show for *before*
 [1...15] are the pages to show
 +, 16 means the *after* link goes to page 16
 if we're on page 65
 that falls within the 61-75 range
 the *before* would be 46
 the *next would be 76
 [ [ '+', 46 ],
  [ 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75 ],
  [ '+', 76 ] ]
 if we're on page 92 of 94
 that falls with the 91-105 range (but we don't have 105 pages)
 so *before* would be 76
 *next* would be no page
 and [91...94] are the pages to show
[ [ '+', 76 ], [ 91, 92, 93, 94 ], [ '-' ] ]
 *
 */
//const PAGE_BY = 15
//import PAGE_BY from './actions/types'
import * as constants from './actions/types'


function pageArrays(totalPages, currentPage) {

  let returnArray = []

  if (totalPages <= constants.PAGE_BY) {
    let pageArray = _.range(1, totalPages + 1)
    returnArray.push(['-'])
    returnArray.push(pageArray)
    returnArray.push(['-'])
    return returnArray
  }
  
  let partitions = Math.floor(totalPages/constants.PAGE_BY) 

  // which segment are we in ??
  let currentPartition = Math.floor(currentPage/constants.PAGE_BY)
  
  let isEnd =  currentPage % constants.PAGE_BY == 0
  if (isEnd) {
    // if it's exact, we don't need to switch to next range
    currentPartition = currentPartition - 1
  }

  let start = (currentPartition * constants.PAGE_BY) + 1
  let end = (start + constants.PAGE_BY > totalPages) ? totalPages : (start + constants.PAGE_BY)
  let pageRange = _.range(start, end)

  if (currentPartition >= partitions) {
    returnArray.push(['+', (currentPartition - 1) * constants.PAGE_BY + 1])
    returnArray.push(pageRange)
    returnArray.push(['-'])
  } else if ((currentPartition < partitions) && (currentPartition > 1)) {
    returnArray.push(['+', (currentPartition - 1) * constants.PAGE_BY + 1])
    returnArray.push(pageRange)
    returnArray.push(['+', ((currentPartition + 1) * constants.PAGE_BY) + 1])
  } else if (currentPartition == 1) {
    returnArray.push(['+', 1])
    returnArray.push(pageRange)
    returnArray.push(['+', ((currentPartition + 1) * constants.PAGE_BY) + 1])
  } else if (currentPartition == 0) {
    returnArray.push(['-'])
    returnArray.push(pageRange)
    returnArray.push(['+', ((currentPartition + 1) * constants.PAGE_BY) + 1])
  }

  return returnArray
}


export default { pageArrays }

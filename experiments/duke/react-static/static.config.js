import axios from 'axios'
import Scholars from './src/fetchers/scholars'
import { makePageRoutes } from 'react-static/node'

export default {
  getSiteData: () => ({
    title: 'Scholars Graphql Example',
  }),
  getRoutes: async () => {
    const { data: people } = await Scholars.allPeople()
    return [
      {
        path: '/',
        component: 'src/containers/Home',
      },
      ...makePageRoutes({
        items: people.personList,
        pageSize: 10,
        pageToken: 'page',
        route: {
          path: '/people',
          component: 'src/containers/people',
        },
        decorate: (items, i, totalPages) => ({
          getData: () => ({
            people: items,
            currentPage: i,
            totalPages
          })
        })
      }),
      {
        path: '/about',
        component: 'src/containers/About',
      },
      {
        is404: true,
        component: 'src/containers/404',
      },
    ]
  },
  // webpack: (config, { defaultLoaders }) => {
  //   config.module.rules = [
  //     {
  //       oneOf: [
  //         {
  //           test: /\.json$/,
  //           use: [{ loader: 'json-loader' }],
  //         },
  //         defaultLoaders.jsLoader,
  //         defaultLoaders.cssLoader,
  //         defaultLoaders.fileLoader,
  //       ],
  //     },
  //   ]
  //   return config
  // },
}

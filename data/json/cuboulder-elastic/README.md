Json files generated from a search against the CU Boulder elastic endpoint.

These elastic structures were built to support the facetview code from cottageLabs.
https://github.com/CottageLabs/facetview2
It's important to note that Elastic and the Elastic json structures were implemented to support the requirements of the facetview javascript library, not the other way around. Meaning that the UI requirements dictated the solution.

CU Elastic endpoint structures can be discover via the URLs:
Publiations: https://experts.colorado.edu/es/fispubs-staging-v1/_search  supporting: https://experts.colorado.edu/publications
People: https://experts.colorado.edu/es/fispeople-staging-v1/_search
Mappings: https://experts.colorado.edu/es/_mapping

Querying syntax for Elastic is easily findable via the web.


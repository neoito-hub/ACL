/* eslint-disable jsx-a11y/control-has-associated-label */
/* eslint-disable camelcase */
/* eslint-disable react/prop-types */
import React from 'react'
import Pagination from '../../Layout/pagination/pagination'

const AddNewPolicy = (props) => {
  const page_limit = Number(process.env.PAGE_LIMIT)
  const {
    // policies,
    loader,
    policiesLoader,
    policiesToAdd,
    policiesSelectedPage,
    policiesTotalCount,
    handlePageChange,
    onSearchTextChange,
    onPolicyChange,
    // selectedPolicies,
  } = props
  return (
    <div className="fadeIn clear-both mt-20 w-full">
      <p className="text-ab-sm pb-3 font-semibold">Add New Policies</p>
      <input
        type="text"
        onChange={onSearchTextChange}
        className="search-input border-ab-gray-dark text-ab-sm h-10 w-full rounded-md !bg-[length:14px_14px] px-2 pl-9 focus:outline-none"
        placeholder="Search for Policies"
      />
      <div className="border-ab-gray-dark custom-h-scroll-primary mt-3.5 w-full overflow-x-auto border">
        <table className="text-ab-black w-full text-left">
          <thead>
            <tr className="bg-ab-gray-light">
              <th className="text-ab-sm whitespace-nowrap p-3 font-medium">
                Policy
              </th>
              <th className="text-ab-sm whitespace-nowrap p-3 font-medium">
                Type
              </th>
              <th className="text-ab-sm whitespace-nowrap p-3 text-center font-medium">
                Attach Policy
              </th>
            </tr>
          </thead>
          <tbody>
            {!loader &&
              policiesToAdd?.map((policy) => (
                <tr
                  key={policy.ac_pol_grp_id}
                  className="text-ab-sm border-ab-gray-dark border-t"
                >
                  <td className="p-3 text-xs">{policy.name}</td>
                  <td className="p-3 text-xs">
                    {policy.is_predefined ? 'Predefined' : 'Custom'}
                  </td>
                  <td className="p-3">
                    <label className="flex justify-center">
                      <input
                        onChange={(e) => onPolicyChange(e, policy)}
                        disabled={policiesLoader}
                        checked={!!policy.subs_id}
                        className="peer hidden"
                        type="checkbox"
                      />
                      <span className="chkbox-icon border-ab-disabled float-left h-5 w-5 flex-shrink-0 cursor-pointer rounded border bg-white" />
                    </label>
                  </td>
                </tr>
              ))}
            {!loader && !policiesToAdd && (
              <tr className="flex justify-center items-center">
                <td className="text-ab-black float-left w-full py-10 text-center text-sm">
                  No Policies Found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <Pagination
        Padding={0}
        marginTop={1}
        pageCount={Math.ceil(policiesTotalCount / page_limit)}
        handlePageChange={handlePageChange}
        selected={policiesSelectedPage}
      />
    </div>
  )
}
export default AddNewPolicy

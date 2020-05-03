import { gql } from 'apollo-boost';

export const USERS_QUERY = gql`
  query listUsers($search: String, $offset: Int, $limit: Int) {
    users(search: $search, offset: $offset, limit: $limit) {
      users {
        id
        context
        updatedAt
      }
      total
    }
  }
`;

// export const TOGGLE_FLAG_QUERY = gql`
//   mutation toggleUser($id: ID!, $input: UpdateUser!) {
//     updateUser(id: $id, input: $input) {
//       id
//       enabled
//     }
//   }
// `;

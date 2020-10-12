// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Define the session user
export class SessionUser {
    user_id: number;
    username: string;
    email: string;
    realname: string;
    role_name?: string;
    role_id?: number;
    has_admin_role?: boolean;
    comment: string;
    oidc_user_meta?: OidcUserMeta;
}
export class OidcUserMeta {
    id: number;
    user_id: number;
    secret: string;
    subiss: string;
    creation_time: Date;
    update_time: Date;
}

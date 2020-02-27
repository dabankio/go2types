/* tslint:disable */
/* eslint-disable */

export namespace types {

  //github.com/dabankio/go2typings/example/types.WeekDay
  export type WeekDay = "mon" | "sun";

  //github.com/dabankio/go2typings/example/types.WeekDay2
  export type WeekDay2 = "3" | "4";

  //github.com/dabankio/go2typings/example/types.T
  export interface T {
    weekday: WeekDay;
    weekday2: WeekDay2;
  }

  //github.com/dabankio/go2typings/example/types.UserTag
  export interface UserTag {
    tag: string;
  }

  //github.com/dabankio/go2typings/example/types.User
  export interface User {
    firstname: string;
    secondName: string;
    tags: Array<UserTag> | null;
  }

  //github.com/dabankio/go2typings/example/user.XTime
  export interface XTime {
  }

  //github.com/dabankio/go2typings/example/user.PersonEstate
  export interface PersonEstate {
    Amount: number;
    Time: number;
  }

  //github.com/dabankio/go2typings/example/user.Person
  //some document for type person
  export interface Person {
    Name: string | null;
    nickname?: string;
    Age: number;//some document for field age
    estates: Array<PersonEstate> | null;
  }

}

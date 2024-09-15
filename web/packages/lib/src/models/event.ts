

export type ContenderPublicInfoUpdatedEvent = {
    contenderId: number,
    compClassId: number,
    publicName: string,
    clubName: string,
    withdrawnFromFinals: boolean,
    disqualified: boolean,
}

export type ContenderScoreUpdatedEvent = {
    timestamp: string,
    contenderId: number,
    score: number,
    placement: number,
}
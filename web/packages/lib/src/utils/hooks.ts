import { differenceInMilliseconds, isBefore } from "date-fns";
import { writable } from "svelte/store";

export type ContestState = "NOT_STARTED" | "RUNNING" | "GRACE_PERIOD" | "ENDED";

export const useContestState = () => {
    const stateStore = writable<ContestState>("NOT_STARTED");
    const state: {
        intervalTimerId: number;
        startTime: Date;
        endTime: Date;
        gracePeriodEndTime?: Date
    } = {
        intervalTimerId: NaN,
        startTime: new Date(),
        endTime: new Date(),
    };

    const computeState = () => {
        const now = new Date();
        let durationUntilNextState: number = NaN;

        switch (true) {
            case isBefore(now, state.startTime):
                stateStore.set("NOT_STARTED");
                durationUntilNextState = differenceInMilliseconds(state.startTime, now);

                break;
            case isBefore(now, state.endTime):
                stateStore.set("RUNNING");
                durationUntilNextState = differenceInMilliseconds(state.endTime, now);

                break;
            case state.gracePeriodEndTime && isBefore(now, state.gracePeriodEndTime):
                stateStore.set("GRACE_PERIOD");
                durationUntilNextState = differenceInMilliseconds(state.gracePeriodEndTime, now);

                break;
            default:
                stateStore.set("ENDED");
        }


        if (durationUntilNextState) {
            state.intervalTimerId = setTimeout(
                computeState,
                durationUntilNextState,
            );
        }
    }

    computeState()

    return {
        state: stateStore,
        stop: () => {
            clearTimeout(state.intervalTimerId)
        },
        update: (startTime: Date, endTime: Date, gracePeriodEndTime?: Date) => {
            state.startTime = startTime;
            state.endTime = endTime;
            state.gracePeriodEndTime = gracePeriodEndTime;

            clearInterval(state.intervalTimerId);
            computeState();
        }
    }
}
import { createQuery } from "@tanstack/svelte-query";
import { ApiClient } from "../Api";

export const getSelfQuery = () =>
    createQuery({
        queryKey: ["self"],
        queryFn: async () => ApiClient.getInstance().getSelf(),
    });
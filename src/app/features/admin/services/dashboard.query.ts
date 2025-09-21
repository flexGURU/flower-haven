import { inject } from "@angular/core";
import { DashboardService } from "./dashboard";
import { injectQuery } from "@tanstack/angular-query-experimental";
import { lastValueFrom } from "rxjs";

export const dashboardQuery = () => {
    const dashboardService = inject(DashboardService);

    const query = injectQuery(() => ({
        queryKey: ['dashboard', dashboardService.getDashboardStats()],
        queryFn: () => lastValueFrom(dashboardService.getDashboardStats()),
        staleTime: 200 * 1000, // 
        refetchOnWindowFocus: true,
    }));

    return query;
};

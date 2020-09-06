import format from 'date-fns/format';
export function localeDate(date: string): string {
    return format(new Date(date), 'ee-MM-yyyy');
}

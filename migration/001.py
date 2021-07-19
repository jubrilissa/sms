import psycopg2
import csv

connection = psycopg2.connect("dbname=dev_sms_thirdterm_backup user=masterp host=localhost password=''")



def get_first_student_scores():
    with open('student_subject_classes.csv')as csvfile:
        score_list = []
        reader = csv.DictReader(csvfile, delimiter=",")
        for current_reader in reader:
            score_list.append(current_reader)
        return score_list


def get_second_student_scores():
    with open('student_subject_classes_2.csv')as csvfile:
        score_list = []
        reader = csv.DictReader(csvfile, delimiter=",")
        for current_reader in reader:
            score_list.append(current_reader)
        return score_list


def update_first_term_student_scores():
    first_student_scores = get_first_student_scores()
    print(f'Total student scores {len(first_student_scores)}')
    with connection.cursor() as cursor:
        updated_row = 0
        for first_student_score in first_student_scores:
            first_ca = first_student_score['s_first_ca'] or None
            second_ca = first_student_score['s_second_ca'] or None
            first_exam = first_student_score['second_exam'] or None
            total_first = first_student_score['total_second'] or None
            student_id = first_student_score['student_id']
            subject_class_id = first_student_score['subject_class_id']
            print(first_student_score)

            update_query = "UPDATE student_subject_classes SET first_ca=(%s), second_ca=(%s), first_exam=(%s), total_first=(%s) WHERE student_id=(%s) AND subject_class_id=(%s)"
            cursor.execute(update_query, (first_ca, second_ca, first_exam, total_first, student_id, subject_class_id))
            print(cursor.rowcount)
            updated_row = updated_row + cursor.rowcount
            print(f'Current row updated {updated_row}')
            connection.commit()
    print(f'Total student scores {len(first_student_scores)}')


def update_second_term_student_scores():
    first_student_scores = get_second_student_scores()
    print(f'Total student scores {len(first_student_scores)}')
    with connection.cursor() as cursor:
        updated_row = 0
        for first_student_score in first_student_scores:
            s_first_ca = first_student_score['first_ca'] or None
            s_second_ca = first_student_score['second_ca'] or None
            second_exam = first_student_score['first_exam'] or None
            total_second = first_student_score['total_first'] or None
            student_id = first_student_score['student_id']
            subject_class_id = first_student_score['subject_class_id']
            print(first_student_score)

            update_query = "UPDATE student_subject_classes SET s_first_ca=(%s), s_second_ca=(%s), second_exam=(%s), total_second=(%s) WHERE student_id=(%s) AND subject_class_id=(%s)"
            cursor.execute(update_query, (s_first_ca, s_second_ca, second_exam, total_second, student_id, subject_class_id))
            print(cursor.rowcount)
            updated_row = updated_row + cursor.rowcount
            print(f'Current row updated {updated_row}')
            connection.commit()
    print(f'Total student scores {len(first_student_scores)}')

# update_first_term_student_scores()

update_second_term_student_scores()